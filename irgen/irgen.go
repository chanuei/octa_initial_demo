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
			ptrType := C.LLVMTypeOf(ptr)
			val := C.LLVMBuildLoad2(ir.builder, ptrType, ptr, C.CString(s.Expr.(ast.VarExpr).Name))
			ir.printInt(val)
		}
	}

	C.LLVMBuildRet(ir.builder, C.LLVMConstInt(C.LLVMInt32Type(), 0, 0))
}

func (ir *IRGen) printInt(val C.LLVMValueRef) {
	fmt.Println("Printing integer (LLVM IR value):", val)
}

// WriteObject 写出 .oo 文件
func (ir *IRGen) WriteObject(filename string) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	if C.LLVMWriteBitcodeToFile(ir.mod, cfilename) != 0 {
		panic("failed to write bitcode file")
	}
}

// Link 链接多个 .oo 文件生成可执行程序
func Link(objects []string, exeName string) error {
	args := append(objects, "-o", exeName)
	fmt.Println("Linking using clang:", args)
	cmd := exec.Command("clang", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
