/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: srv.proto */

#ifndef PROTOBUF_C_srv_2eproto__INCLUDED
#define PROTOBUF_C_srv_2eproto__INCLUDED

#include <protobuf-c/protobuf-c.h>

PROTOBUF_C__BEGIN_DECLS

#if PROTOBUF_C_VERSION_NUMBER < 1003000
# error This file was generated by a newer version of protoc-c which is incompatible with your libprotobuf-c headers. Please update your headers.
#elif 1003001 < PROTOBUF_C_MIN_COMPILER_VERSION
# error This file was generated by an older version of protoc-c which is incompatible with your libprotobuf-c headers. Please regenerate this file with a newer version of protoc-c.
#endif


typedef struct _Srv__NotifyReadyReq Srv__NotifyReadyReq;
typedef struct _Srv__BioErrorReq Srv__BioErrorReq;


/* --- enums --- */


/* --- messages --- */

struct  _Srv__NotifyReadyReq
{
  ProtobufCMessage base;
  /*
   * CaRT URI
   */
  char *uri;
  /*
   * Number of CaRT contexts
   */
  uint32_t nctxs;
  /*
   * Path to IO server's dRPC listener socket
   */
  char *drpclistenersock;
  /*
   * IO server instance index
   */
  uint32_t instanceidx;
};
#define SRV__NOTIFY_READY_REQ__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&srv__notify_ready_req__descriptor) \
    , (char *)protobuf_c_empty_string, 0, (char *)protobuf_c_empty_string, 0 }


struct  _Srv__BioErrorReq
{
  ProtobufCMessage base;
  /*
   * unmap I/O error
   */
  protobuf_c_boolean unmaperr;
  /*
   * read I/O error
   */
  protobuf_c_boolean readerr;
  /*
   * write I/O error
   */
  protobuf_c_boolean writeerr;
  /*
   * VOS target ID
   */
  int32_t tgtid;
  /*
   * IO server instance index
   */
  uint32_t instanceidx;
  /*
   * Path to IO server's dRPC listener socket
   */
  char *drpclistenersock;
  /*
   * CaRT URI
   */
  char *uri;
};
#define SRV__BIO_ERROR_REQ__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&srv__bio_error_req__descriptor) \
    , 0, 0, 0, 0, 0, (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string }


/* Srv__NotifyReadyReq methods */
void   srv__notify_ready_req__init
                     (Srv__NotifyReadyReq         *message);
size_t srv__notify_ready_req__get_packed_size
                     (const Srv__NotifyReadyReq   *message);
size_t srv__notify_ready_req__pack
                     (const Srv__NotifyReadyReq   *message,
                      uint8_t             *out);
size_t srv__notify_ready_req__pack_to_buffer
                     (const Srv__NotifyReadyReq   *message,
                      ProtobufCBuffer     *buffer);
Srv__NotifyReadyReq *
       srv__notify_ready_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   srv__notify_ready_req__free_unpacked
                     (Srv__NotifyReadyReq *message,
                      ProtobufCAllocator *allocator);
/* Srv__BioErrorReq methods */
void   srv__bio_error_req__init
                     (Srv__BioErrorReq         *message);
size_t srv__bio_error_req__get_packed_size
                     (const Srv__BioErrorReq   *message);
size_t srv__bio_error_req__pack
                     (const Srv__BioErrorReq   *message,
                      uint8_t             *out);
size_t srv__bio_error_req__pack_to_buffer
                     (const Srv__BioErrorReq   *message,
                      ProtobufCBuffer     *buffer);
Srv__BioErrorReq *
       srv__bio_error_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   srv__bio_error_req__free_unpacked
                     (Srv__BioErrorReq *message,
                      ProtobufCAllocator *allocator);
/* --- per-message closures --- */

typedef void (*Srv__NotifyReadyReq_Closure)
                 (const Srv__NotifyReadyReq *message,
                  void *closure_data);
typedef void (*Srv__BioErrorReq_Closure)
                 (const Srv__BioErrorReq *message,
                  void *closure_data);

/* --- services --- */


/* --- descriptors --- */

extern const ProtobufCMessageDescriptor srv__notify_ready_req__descriptor;
extern const ProtobufCMessageDescriptor srv__bio_error_req__descriptor;

PROTOBUF_C__END_DECLS


#endif  /* PROTOBUF_C_srv_2eproto__INCLUDED */
