// @generated by protoc-gen-connect-es v0.13.0 with parameter "target=ts"
// @generated from file golink/v1/debug.proto (package golink.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { DebugRequest, DebugResponse } from "./debug_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service golink.v1.DebugService
 */
export const DebugService = {
  typeName: "golink.v1.DebugService",
  methods: {
    /**
     * @generated from rpc golink.v1.DebugService.Debug
     */
    debug: {
      name: "Debug",
      I: DebugRequest,
      O: DebugResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

