from concurrent import futures
import logging


import grpc
from grpc_reflection.v1alpha import reflection
import service.recognizer_pb2_grpc as recognizer_pb2_grpc
import service.recognizer_pb2 as recognizer_pb2


from latex_ocr import LatexOCR

class RecognizerServicer(recognizer_pb2_grpc.RecognizerServicer):
    def __init__(self):
        self.latex_ocr = LatexOCR()

    def RecognizeFormula(self, request, context):
        formula = self.latex_ocr.convert_formula(request.data[0])
        return recognizer_pb2.RecognizedResult(result=formula)

    def RecognizeMixed(self, request, context):
        result = ''
        for img in request.data:
            result += self.latex_ocr.convert_mixed(img) + '\n'

        return recognizer_pb2.RecognizedResult(result=result)

    def RecognizeText(self, request, context):
        result = ''
        for img in request.data:
            result += self.latex_ocr.convert_text(img) + '\n'
        return recognizer_pb2.RecognizedResult(result=result)
