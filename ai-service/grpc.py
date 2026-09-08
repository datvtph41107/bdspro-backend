# """
# gRPC Server cho OCR National Card Service
# Port: 8219
# """

# import os
# import sys
# import time
# from concurrent import futures
# import grpc

# # Import generated protobuf code
# try:
#     import ocr_national_card_pb2
#     import ocr_national_card_pb2_grpc
# except ImportError:
#     print("⚠️  Warning: Protobuf files not generated.")
#     print("   Run: python -m grpc_tools.protoc --python_out=. --grpc_python_out=. --proto_path=. ocr_national_card.proto")
#     sys.exit(1)

# # Import OCR module
# from app.ocr.national_card_ocr import NationalCardOCR


# class OcrServiceServicer(ocr_national_card_pb2_grpc.OcrServiceServicer):
#     """gRPC service implementation cho OCR National Card"""
    
#     def __init__(self):
#         """Khởi tạo OCR engine"""
#         print("🔧 Initializing OCR engine...")
#         try:
#             self.ocr_engine = NationalCardOCR()
#             print("✅ OCR engine initialized successfully")
#         except Exception as e:
#             print(f"❌ Failed to initialize OCR engine: {e}")
#             raise
    
#     def OcrNationalCard(
#         self, 
#         request: ocr_national_card_pb2.OcrNationalCardRequest, 
#         context: grpc.ServicerContext
#     ) -> ocr_national_card_pb2.OcrNationalCardResponse:
#         """
#         RPC để trích xuất thông tin từ ảnh căn cước
        
#         Args:
#             request: OcrNationalCardRequest chứa front_image_url và back_image_url
#             context: gRPC context
            
#         Returns:
#             OcrNationalCardResponse chứa thông tin đã trích xuất
#         """
#         try:
#             # Validate request
#             if not request.front_image_url or not request.back_image_url:
#                 context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
#                 context.set_details("front_image_url and back_image_url are required")
#                 return ocr_national_card_pb2.OcrNationalCardResponse()
            
#             print(f"📸 Processing national card images:")
#             print(f"   Front: {request.front_image_url}")
#             print(f"   Back: {request.back_image_url}")
            
#             # Extract information
#             info_dict, mrz, processing_time_ms = self.ocr_engine.extract_national_card_info(
#                 front_image_url=request.front_image_url,
#                 back_image_url=request.back_image_url
#             )
            
#             # Build response
#             card_info = ocr_national_card_pb2.NationalCardInfo()
            
#             # Set fields from info_dict
#             if 'id_number' in info_dict and info_dict['id_number']:
#                 card_info.id_number = info_dict['id_number']
#             if 'full_name' in info_dict and info_dict['full_name']:
#                 card_info.full_name = info_dict['full_name']
#             if 'date_of_birth' in info_dict and info_dict['date_of_birth']:
#                 card_info.date_of_birth = info_dict['date_of_birth']
#             if 'gender' in info_dict and info_dict['gender']:
#                 card_info.gender = info_dict['gender']
#             if 'nationality' in info_dict and info_dict['nationality']:
#                 card_info.nationality = info_dict['nationality']
#             if 'place_of_origin' in info_dict and info_dict['place_of_origin']:
#                 card_info.place_of_origin = info_dict['place_of_origin']
#             if 'place_of_residence' in info_dict and info_dict['place_of_residence']:
#                 card_info.place_of_residence = info_dict['place_of_residence']
#             if 'ethnic' in info_dict and info_dict['ethnic']:
#                 card_info.ethnic = info_dict['ethnic']
#             if 'religion' in info_dict and info_dict['religion']:
#                 card_info.religion = info_dict['religion']
#             if 'issue_date' in info_dict and info_dict['issue_date']:
#                 card_info.issue_date = info_dict['issue_date']
#             if 'issue_place' in info_dict and info_dict['issue_place']:
#                 card_info.issue_place = info_dict['issue_place']
#             if 'expiry_date' in info_dict and info_dict['expiry_date']:
#                 card_info.expiry_date = info_dict['expiry_date']
#             if 'features' in info_dict and info_dict['features']:
#                 card_info.features = info_dict['features']
#             if 'address' in info_dict and info_dict['address']:
#                 card_info.address = info_dict['address']
            
#             # Build metadata
#             metadata = ocr_national_card_pb2.ExtractMetadata(
#                 processing_time_ms=int(processing_time_ms),
#                 model="tesseract-ocr",
#                 confidence_score=0.85  # TODO: Calculate actual confidence
#             )
            
#             # Build response
#             response = ocr_national_card_pb2.OcrNationalCardResponse(
#                 info=card_info,
#                 mrz=mrz or "",
#                 metadata=metadata
#             )
            
#             print(f"✅ Successfully extracted information (took {int(processing_time_ms)}ms)")
            
#             return response
            
#         except ValueError as e:
#             # Invalid input
#             context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
#             context.set_details(str(e))
#             return ocr_national_card_pb2.OcrNationalCardResponse()
            
#         except Exception as e:
#             # Internal error
#             import traceback
#             traceback.print_exc()
#             context.set_code(grpc.StatusCode.INTERNAL)
#             context.set_details(f"OCR processing failed: {str(e)}")
#             return ocr_national_card_pb2.OcrNationalCardResponse()


# def serve():
#     """Khởi động gRPC server"""
#     port = int(os.getenv("GRPC_PORT", 8219))
    
#     # Tạo gRPC server
#     server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
#     # Đăng ký service
#     ocr_national_card_pb2_grpc.add_OcrServiceServicer_to_server(
#         OcrServiceServicer(), 
#         server
#     )
    
#     # Bind port
#     server.add_insecure_port(f'[::]:{port}')
    
#     # Start server
#     server.start()
    
#     print(f"🚀 gRPC OCR Service started on port {port}")
#     print(f"   Service: OcrService")
#     print(f"   RPC: OcrNationalCard")
#     print(f"   Waiting for requests...")
    
#     try:
#         server.wait_for_termination()
#     except KeyboardInterrupt:
#         print("\n🛑 Shutting down gRPC server...")
#         server.stop(0)
#         print("✅ Server stopped")


# if __name__ == "__main__":
#     serve()

