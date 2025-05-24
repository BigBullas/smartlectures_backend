from abc import ABC, abstractmethod
from typing import Union


class ABSTextRecognizer(ABC):
    def __init__(self, model):
        pass

    @abstractmethod
    def recognize(self, path: str):
        pass
           