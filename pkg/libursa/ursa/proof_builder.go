package ursa

/*
   #cgo LDFLAGS: -lursa
   #include "ursa_cl.h"
   #include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type ProofBuilder Handle
type ProofHandle Handle

// NewProofBuilder
func NewProofBuilder() (*ProofBuilder, error) {
	var builder unsafe.Pointer

	result := C.ursa_cl_prover_new_proof_builder(&builder)
	if result.code != 0 {
		return nil, ursaError(C.GoString(result.message))
	}

	return &ProofBuilder{builder}, nil
}

// ProofFromJSON creates and returns proof from JSON
func ProofFromJSON(jsn []byte) (*ProofHandle, error) {
	var builder unsafe.Pointer
	cjson := C.CString(string(jsn))
	defer C.free(unsafe.Pointer(cjson))

	result := C.ursa_cl_proof_from_json(cjson, &builder)
	if result.code != 0 {
		return nil, ursaError(C.GoString(result.message))
	}

	return &ProofHandle{builder}, nil
}

// AddCommonAttribute add a common attribute to the proof builder
func (r *ProofBuilder) AddCommonAttribute(attr string) error {
	cattr := C.CString(attr)
	defer C.free(unsafe.Pointer(cattr))

	result := C.ursa_cl_proof_builder_add_common_attribute(r.ptr, cattr)
	if result.code != 0 {
		return ursaError(C.GoString(result.message))
	}

	return nil
}

// AddSubProofRequest add a sub proof request to the proof builder
func (r *ProofBuilder) AddSubProofRequest(subProof *SubProofRequestHandle, credSchema *CredentialSchemaHandle,
	nonCredSchema *NonCredentialSchemaHandle, signature *CredentialSignature, values *CredentialValues, pubKey *CredentialDefPubKey) error {

	result := C.ursa_cl_proof_builder_add_sub_proof_request(r.ptr, subProof.ptr, credSchema.ptr, nonCredSchema.ptr,
		signature.ptr, values.ptr, pubKey.ptr /*revoc_reg*/, nil /*witness*/, nil)
	if result.code != 0 {
		return ursaError(C.GoString(result.message))
	}

	return nil
}

// Finalize finalize the proof
func (r *ProofBuilder) Finalize(nonce *Nonce) (*ProofHandle, error) {
	var proof unsafe.Pointer

	result := C.ursa_cl_proof_builder_finalize(r.ptr, nonce.ptr, &proof)
	if result.code != 0 {
		return nil, ursaError(C.GoString(result.message))
	}

	return &ProofHandle{proof}, nil
}

// ToJSON returns JSON representation of proof
func (r *ProofHandle) ToJSON() ([]byte, error) {
	var d *C.char
	defer C.free(unsafe.Pointer(d))

	result := C.ursa_cl_proof_to_json(r.ptr, &d)
	if result.code != 0 {
		return nil, ursaError(C.GoString(result.message))
	}

	return []byte(C.GoString(d)), nil
}

// Free deallocates proof instance
func (r *ProofHandle) Free() error {
	result := C.ursa_cl_proof_free(r.ptr)
	if result.code != 0 {
		return ursaError(C.GoString(result.message))
	}

	return nil
}

type SubProofRequestBuilder Handle
type SubProofRequestHandle Handle

// NewSubProofRequestBuilder creates and returns sub proof request entity builder
func NewSubProofRequestBuilder() (*SubProofRequestBuilder, error) {
	var builder unsafe.Pointer

	result := C.ursa_cl_sub_proof_request_builder_new(&builder)
	if result.code != 0 {
		return nil, ursaError(C.GoString(result.message))
	}

	return &SubProofRequestBuilder{builder}, nil

}

// AddPredicate adds predicate to sub proof request
func (r *SubProofRequestBuilder) AddPredicate(attr, ptype string, value int32) error {
	cattr := C.CString(attr)
	defer C.free(unsafe.Pointer(cattr))
	cptype := C.CString(ptype)
	defer C.free(unsafe.Pointer(cptype))

	result := C.ursa_cl_sub_proof_request_builder_add_predicate(r.ptr, cattr, cptype, C.int32_t(value))
	if result.code != 0 {
		return ursaError(C.GoString(result.message))
	}

	return nil
}

// AddRevealedAttr adds new revealed attribute to sub proof request.
func (r *SubProofRequestBuilder) AddRevealedAttr(attr string) error {
	cattr := C.CString(attr)
	defer C.free(unsafe.Pointer(cattr))

	result := C.ursa_cl_sub_proof_request_builder_add_revealed_attr(r.ptr, cattr)
	if result.code != 0 {
		return ursaError(C.GoString(result.message))
	}

	return nil
}

// Finalize deallocates sub proof request builder and returns sub proof request entity instead.
func (r *SubProofRequestBuilder) Finalize() (*SubProofRequestHandle, error) {
	var proof unsafe.Pointer

	result := C.ursa_cl_sub_proof_request_builder_finalize(r.ptr, &proof)
	if result.code != 0 {
		return nil, ursaError(C.GoString(result.message))
	}

	return &SubProofRequestHandle{proof}, nil
}

// Free deallocates sub proof request instance
func (r *SubProofRequestHandle) Free() error {
	result := C.ursa_cl_sub_proof_request_free(r.ptr)
	if result.code != 0 {
		return ursaError(C.GoString(result.message))
	}

	return nil
}
