/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: mgmt.proto */

#ifndef PROTOBUF_C_mgmt_2eproto__INCLUDED
#define PROTOBUF_C_mgmt_2eproto__INCLUDED

#include <protobuf-c/protobuf-c.h>

PROTOBUF_C__BEGIN_DECLS

#if PROTOBUF_C_VERSION_NUMBER < 1003000
# error This file was generated by a newer version of protoc-c which is incompatible with your libprotobuf-c headers. Please update your headers.
#elif 1003002 < PROTOBUF_C_MIN_COMPILER_VERSION
# error This file was generated by an older version of protoc-c which is incompatible with your libprotobuf-c headers. Please regenerate this file with a newer version of protoc-c.
#endif

#include "pool.pb-c.h"
#include "storage_query.pb-c.h"

typedef struct _Mgmt__JoinReq Mgmt__JoinReq;
typedef struct _Mgmt__JoinResp Mgmt__JoinResp;
typedef struct _Mgmt__GetAttachInfoReq Mgmt__GetAttachInfoReq;
typedef struct _Mgmt__GetAttachInfoResp Mgmt__GetAttachInfoResp;
typedef struct _Mgmt__GetAttachInfoResp__Psr Mgmt__GetAttachInfoResp__Psr;


/* --- enums --- */

/*
 * Server state in the system map.
 */
typedef enum _Mgmt__JoinResp__State {
  /*
   * Server in the system.
   */
  MGMT__JOIN_RESP__STATE__IN = 0,
  /*
   * Server excluded from the system.
   */
  MGMT__JOIN_RESP__STATE__OUT = 1
    PROTOBUF_C__FORCE_ENUM_TO_BE_INT_SIZE(MGMT__JOIN_RESP__STATE)
} Mgmt__JoinResp__State;

/* --- messages --- */

struct  _Mgmt__JoinReq
{
  ProtobufCMessage base;
  /*
   * Server UUID.
   */
  char *uuid;
  /*
   * Server rank desired, if not -1.
   */
  uint32_t rank;
  /*
   * Server CaRT base URI (i.e., for context 0).
   */
  char *uri;
  /*
   * Server CaRT context count.
   */
  uint32_t nctxs;
  /*
   * Server management address.
   */
  char *addr;
};
#define MGMT__JOIN_REQ__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&mgmt__join_req__descriptor) \
    , (char *)protobuf_c_empty_string, 0, (char *)protobuf_c_empty_string, 0, (char *)protobuf_c_empty_string }


struct  _Mgmt__JoinResp
{
  ProtobufCMessage base;
  /*
   * DAOS error code
   */
  int32_t status;
  /*
   * Server rank assigned.
   */
  uint32_t rank;
  Mgmt__JoinResp__State state;
};
#define MGMT__JOIN_RESP__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&mgmt__join_resp__descriptor) \
    , 0, 0, MGMT__JOIN_RESP__STATE__IN }


struct  _Mgmt__GetAttachInfoReq
{
  ProtobufCMessage base;
  /*
   * System name. For daos_agent only.
   */
  char *sys;
};
#define MGMT__GET_ATTACH_INFO_REQ__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&mgmt__get_attach_info_req__descriptor) \
    , (char *)protobuf_c_empty_string }


/*
 * CaRT PSR.
 */
struct  _Mgmt__GetAttachInfoResp__Psr
{
  ProtobufCMessage base;
  uint32_t rank;
  char *uri;
};
#define MGMT__GET_ATTACH_INFO_RESP__PSR__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&mgmt__get_attach_info_resp__psr__descriptor) \
    , 0, (char *)protobuf_c_empty_string }


struct  _Mgmt__GetAttachInfoResp
{
  ProtobufCMessage base;
  /*
   * DAOS error code
   */
  int32_t status;
  /*
   * CaRT PSRs of the system group.
   */
  size_t n_psrs;
  Mgmt__GetAttachInfoResp__Psr **psrs;
};
#define MGMT__GET_ATTACH_INFO_RESP__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&mgmt__get_attach_info_resp__descriptor) \
    , 0, 0,NULL }


/* Mgmt__JoinReq methods */
void   mgmt__join_req__init
                     (Mgmt__JoinReq         *message);
size_t mgmt__join_req__get_packed_size
                     (const Mgmt__JoinReq   *message);
size_t mgmt__join_req__pack
                     (const Mgmt__JoinReq   *message,
                      uint8_t             *out);
size_t mgmt__join_req__pack_to_buffer
                     (const Mgmt__JoinReq   *message,
                      ProtobufCBuffer     *buffer);
Mgmt__JoinReq *
       mgmt__join_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   mgmt__join_req__free_unpacked
                     (Mgmt__JoinReq *message,
                      ProtobufCAllocator *allocator);
/* Mgmt__JoinResp methods */
void   mgmt__join_resp__init
                     (Mgmt__JoinResp         *message);
size_t mgmt__join_resp__get_packed_size
                     (const Mgmt__JoinResp   *message);
size_t mgmt__join_resp__pack
                     (const Mgmt__JoinResp   *message,
                      uint8_t             *out);
size_t mgmt__join_resp__pack_to_buffer
                     (const Mgmt__JoinResp   *message,
                      ProtobufCBuffer     *buffer);
Mgmt__JoinResp *
       mgmt__join_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   mgmt__join_resp__free_unpacked
                     (Mgmt__JoinResp *message,
                      ProtobufCAllocator *allocator);
/* Mgmt__GetAttachInfoReq methods */
void   mgmt__get_attach_info_req__init
                     (Mgmt__GetAttachInfoReq         *message);
size_t mgmt__get_attach_info_req__get_packed_size
                     (const Mgmt__GetAttachInfoReq   *message);
size_t mgmt__get_attach_info_req__pack
                     (const Mgmt__GetAttachInfoReq   *message,
                      uint8_t             *out);
size_t mgmt__get_attach_info_req__pack_to_buffer
                     (const Mgmt__GetAttachInfoReq   *message,
                      ProtobufCBuffer     *buffer);
Mgmt__GetAttachInfoReq *
       mgmt__get_attach_info_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   mgmt__get_attach_info_req__free_unpacked
                     (Mgmt__GetAttachInfoReq *message,
                      ProtobufCAllocator *allocator);
/* Mgmt__GetAttachInfoResp__Psr methods */
void   mgmt__get_attach_info_resp__psr__init
                     (Mgmt__GetAttachInfoResp__Psr         *message);
/* Mgmt__GetAttachInfoResp methods */
void   mgmt__get_attach_info_resp__init
                     (Mgmt__GetAttachInfoResp         *message);
size_t mgmt__get_attach_info_resp__get_packed_size
                     (const Mgmt__GetAttachInfoResp   *message);
size_t mgmt__get_attach_info_resp__pack
                     (const Mgmt__GetAttachInfoResp   *message,
                      uint8_t             *out);
size_t mgmt__get_attach_info_resp__pack_to_buffer
                     (const Mgmt__GetAttachInfoResp   *message,
                      ProtobufCBuffer     *buffer);
Mgmt__GetAttachInfoResp *
       mgmt__get_attach_info_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   mgmt__get_attach_info_resp__free_unpacked
                     (Mgmt__GetAttachInfoResp *message,
                      ProtobufCAllocator *allocator);
/* --- per-message closures --- */

typedef void (*Mgmt__JoinReq_Closure)
                 (const Mgmt__JoinReq *message,
                  void *closure_data);
typedef void (*Mgmt__JoinResp_Closure)
                 (const Mgmt__JoinResp *message,
                  void *closure_data);
typedef void (*Mgmt__GetAttachInfoReq_Closure)
                 (const Mgmt__GetAttachInfoReq *message,
                  void *closure_data);
typedef void (*Mgmt__GetAttachInfoResp__Psr_Closure)
                 (const Mgmt__GetAttachInfoResp__Psr *message,
                  void *closure_data);
typedef void (*Mgmt__GetAttachInfoResp_Closure)
                 (const Mgmt__GetAttachInfoResp *message,
                  void *closure_data);

/* --- services --- */

typedef struct _Mgmt__MgmtSvc_Service Mgmt__MgmtSvc_Service;
struct _Mgmt__MgmtSvc_Service
{
  ProtobufCService base;
  void (*join)(Mgmt__MgmtSvc_Service *service,
               const Mgmt__JoinReq *input,
               Mgmt__JoinResp_Closure closure,
               void *closure_data);
  void (*create_pool)(Mgmt__MgmtSvc_Service *service,
                      const Mgmt__CreatePoolReq *input,
                      Mgmt__CreatePoolResp_Closure closure,
                      void *closure_data);
  void (*destroy_pool)(Mgmt__MgmtSvc_Service *service,
                       const Mgmt__DestroyPoolReq *input,
                       Mgmt__DestroyPoolResp_Closure closure,
                       void *closure_data);
  void (*get_attach_info)(Mgmt__MgmtSvc_Service *service,
                          const Mgmt__GetAttachInfoReq *input,
                          Mgmt__GetAttachInfoResp_Closure closure,
                          void *closure_data);
  void (*bio_health_query)(Mgmt__MgmtSvc_Service *service,
                           const Mgmt__BioHealthReq *input,
                           Mgmt__BioHealthResp_Closure closure,
                           void *closure_data);
  void (*smd_list_devs)(Mgmt__MgmtSvc_Service *service,
                        const Mgmt__SmdDevReq *input,
                        Mgmt__SmdDevResp_Closure closure,
                        void *closure_data);
};
typedef void (*Mgmt__MgmtSvc_ServiceDestroy)(Mgmt__MgmtSvc_Service *);
void mgmt__mgmt_svc__init (Mgmt__MgmtSvc_Service *service,
                           Mgmt__MgmtSvc_ServiceDestroy destroy);
#define MGMT__MGMT_SVC__BASE_INIT \
    { &mgmt__mgmt_svc__descriptor, protobuf_c_service_invoke_internal, NULL }
#define MGMT__MGMT_SVC__INIT(function_prefix__) \
    { MGMT__MGMT_SVC__BASE_INIT,\
      function_prefix__ ## join,\
      function_prefix__ ## create_pool,\
      function_prefix__ ## destroy_pool,\
      function_prefix__ ## get_attach_info,\
      function_prefix__ ## bio_health_query,\
      function_prefix__ ## smd_list_devs  }
void mgmt__mgmt_svc__join(ProtobufCService *service,
                          const Mgmt__JoinReq *input,
                          Mgmt__JoinResp_Closure closure,
                          void *closure_data);
void mgmt__mgmt_svc__create_pool(ProtobufCService *service,
                                 const Mgmt__CreatePoolReq *input,
                                 Mgmt__CreatePoolResp_Closure closure,
                                 void *closure_data);
void mgmt__mgmt_svc__destroy_pool(ProtobufCService *service,
                                  const Mgmt__DestroyPoolReq *input,
                                  Mgmt__DestroyPoolResp_Closure closure,
                                  void *closure_data);
void mgmt__mgmt_svc__get_attach_info(ProtobufCService *service,
                                     const Mgmt__GetAttachInfoReq *input,
                                     Mgmt__GetAttachInfoResp_Closure closure,
                                     void *closure_data);
void mgmt__mgmt_svc__bio_health_query(ProtobufCService *service,
                                      const Mgmt__BioHealthReq *input,
                                      Mgmt__BioHealthResp_Closure closure,
                                      void *closure_data);
void mgmt__mgmt_svc__smd_list_devs(ProtobufCService *service,
                                   const Mgmt__SmdDevReq *input,
                                   Mgmt__SmdDevResp_Closure closure,
                                   void *closure_data);

/* --- descriptors --- */

extern const ProtobufCMessageDescriptor mgmt__join_req__descriptor;
extern const ProtobufCMessageDescriptor mgmt__join_resp__descriptor;
extern const ProtobufCEnumDescriptor    mgmt__join_resp__state__descriptor;
extern const ProtobufCMessageDescriptor mgmt__get_attach_info_req__descriptor;
extern const ProtobufCMessageDescriptor mgmt__get_attach_info_resp__descriptor;
extern const ProtobufCMessageDescriptor mgmt__get_attach_info_resp__psr__descriptor;
extern const ProtobufCServiceDescriptor mgmt__mgmt_svc__descriptor;

PROTOBUF_C__END_DECLS


#endif  /* PROTOBUF_C_mgmt_2eproto__INCLUDED */
