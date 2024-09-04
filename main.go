package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const (
	ceSpecVersionHeader  = "ce-specversion"
	cePartitionKeyHeader = "ce-partitionkey"
	ceSourceHeader       = "ce-source"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{contextID: contextID}
}

type pluginContext struct {
	types.DefaultPluginContext
	// In Go, we typically store context-related data directly in structs.
	contextID uint32
}

func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &ceHeaders{contextID: contextID}
}

type ceHeaders struct {
	types.DefaultHttpContext
	contextID uint32
}

var _ types.HttpContext = &ceHeaders{}

func (ctx *ceHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("[ce-modification] on_http_request_headers")

	ceVersion, err := proxywasm.GetHttpRequestHeader(ceSpecVersionHeader)
	if err != nil || ceVersion != "1.0" {
		proxywasm.LogInfof("[ce-modification] invalid ce-specversion header not defined. ce-specversion header value : %s", ceVersion)
		return types.ActionContinue
	}

	if partitionKey, err := proxywasm.GetHttpRequestHeader(cePartitionKeyHeader); err == nil {
		proxywasm.LogInfof("[ce-modification] cloud event partition key set: " + partitionKey)
		return types.ActionContinue
	}

	source, err := proxywasm.GetHttpRequestHeader(ceSourceHeader)
	if err != nil {
		return types.ActionContinue
	}
	proxywasm.LogInfof("[ce-modification] setting cloud event partitionkey: " + source)
	if source != "" {
		if err = proxywasm.AddHttpRequestHeader(cePartitionKeyHeader, source); err != nil {
			proxywasm.LogError("failed to set ce-partitionkey header: " + err.Error())
		}
	}

	// for debug
	h, err := proxywasm.GetHttpRequestHeader(ceSourceHeader)
	proxywasm.LogInfof("%s : %s; err: %v", ceSourceHeader, h, err)

	pk, _ := proxywasm.GetHttpRequestHeader(cePartitionKeyHeader)
	proxywasm.LogInfof("%s : %s, err: %v", cePartitionKeyHeader, pk, err)

	return types.ActionContinue
}
