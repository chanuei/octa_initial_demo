package irgen

/*
#cgo CFLAGS: -I/opt/homebrew/opt/llvm/include
#cgo LDFLAGS: -L/opt/homebrew/opt/llvm/lib -lLLVM
#include <llvm-c/Core.h>
#include <llvm-c/BitWriter.h>
#include <llvm-c/ExecutionEngine.h>
#include <llvm-c/Target.h>
#include <llvm-c/TargetMachine.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"octa/ast"
	"os"
	"os/exec"
	"unsafe"
)

type IRGen struct {
	mod     C.LLVMModuleRef
	builder C.LLVMBuilderRef
}

func NewIRGen() *IRGen {
	mod := C.LLVMModuleCreateWithName(C.CString("octa_module"))
	builder := C.LLVMCreateBuilder()
	return &IRGen{mod: mod, builder: builder}
}

// GenerateFunctions 遍历 AST 生成 LLVM IR
func (ir *IRGen) GenerateFunctions(f *ast.FuncStmt) {
	funcType := C.LLVMFunctionType(C.LLVMInt32Type(), nil, 0, 0)
	funcRef := C.LLVMAddFunction(ir.mod, C.CString(f.Name), funcType)
	C.LLVMSetLinkage(funcRef, C.LLVMLinkage(C.LLVMExternalLinkage)) // ← 使用新版枚举
	entry := C.LLVMAppendBasicBlock(funcRef, C.CString("entry"))
	C.LLVMPositionBuilderAtEnd(ir.builder, entry)

	vars := make(map[string]C.LLVMValueRef)

	for _, stmt := range f.Body {
		switch s := stmt.(type) {
		case *ast.VarDeclStmt:
			ptr := C.LLVMBuildAlloca(ir.builder, C.LLVMInt32Type(), C.CString(s.Name))
			val := C.LLVMConstInt(C.LLVMInt32Type(), C.ulonglong(s.Expr.(ast.NumberExpr).Value), 0)
			C.LLVMBuildStore(ir.builder, val, ptr)
			vars[s.Name] = ptr
		case *ast.AssignStmt:
			ptr, ok := vars[s.Name]
			if !ok {
				panic("variable not declared: " + s.Name)
			}
			val := C.LLVMConstInt(C.LLVMInt32Type(), C.ulonglong(s.Expr.(ast.NumberExpr).Value), 0)
			C.LLVMBuildStore(ir.builder, val, ptr)
		case *ast.PrintStmt:
			ptr, ok := vars[s.Expr.(ast.VarExpr).Name]
			if !ok {
				panic("variable not declared: " + s.Expr.(ast.VarExpr).Name)
			}
			fmt.Println("Preparing to print variable:", s.Expr.(ast.VarExpr).Name)
			ptrType := C.LLVMTypeOf(ptr)
			val := C.LLVMBuildLoad2(ir.builder, ptrType, ptr, C.CString(s.Expr.(ast.VarExpr).Name))
			ir.printInt(val)
		}
	}

	C.LLVMBuildRet(ir.builder, C.LLVMConstInt(C.LLVMInt32Type(), 0, 0))
}

func (ir *IRGen) printInt(val C.LLVMValueRef) {
	fmt.Printf("Printing integer (LLVM IR value): %v\n", val)
	// 先声明 printf
	printfParamTypes := []C.LLVMTypeRef{C.LLVMPointerType(C.LLVMInt8Type(), 0)}
	printfType := C.LLVMFunctionType(
		C.LLVMInt32Type(),
		&printfParamTypes[0],
		1, // 参数个数
		1, // vararg
	)
	printf := C.LLVMGetNamedFunction(ir.mod, C.CString("printf"))
	if printf == nil {
		printf = C.LLVMAddFunction(ir.mod, C.CString("printf"), printfType)
	}

	formatStr := C.LLVMBuildGlobalStringPtr(
		ir.builder,
		C.CString("%d\n"),
		C.CString("fmt"),
	)
	// 构造参数
	args := []C.LLVMValueRef{formatStr, val}

	// 正确调用
	C.LLVMBuildCall2(
		ir.builder,
		printfType, // 直接用声明时保存的函数类型
		printf,
		&args[0],
		C.uint(len(args)),
		C.CString(""),
	)
}

func (ir *IRGen) WriteObject(filename string) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	if C.LLVMWriteBitcodeToFile(ir.mod, cfilename) != 0 {
		panic("failed to write bitcode file")
	}
}

func Link(objects []string, exeName string) error {
	tmpC := "main_wrapper.c"
	cfile, err := os.Create(tmpC)
	if err != nil {
		return err
	}
	defer os.Remove(tmpC)

	cfile.WriteString(`
extern int entrance();
int main() { return entrance(); }
`)
	cfile.Close()

	objFiles := []string{}
	for _, bc := range objects {
		oFile := bc[:len(bc)-3] + ".o"
		cmd := exec.Command("llc", "-filetype=obj", bc, "-o", oFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		objFiles = append(objFiles, oFile)
	}

	args := append(objFiles, tmpC, "-o", exeName)
	cmd := exec.Command("clang", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
