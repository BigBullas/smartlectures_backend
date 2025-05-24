from concurrent import futures
import logging

from service.servicer import RecognizerServicer

import grpc
from grpc_reflection.v1alpha import reflection
import service.recognizer_pb2_grpc as recognizer_pb2_grpc
import service.recognizer_pb2 as recognizer_pb2




def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    recognizer_pb2_grpc.add_RecognizerServicer_to_server(
        RecognizerServicer(), server
    )
    SERVICE_NAMES = (
        recognizer_pb2.DESCRIPTOR.services_by_name['Recognizer'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
