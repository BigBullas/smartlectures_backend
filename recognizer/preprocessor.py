import numpy as np
import cv2
import json
from typing import Union
import scipy

class Preprocessor():
    def __init__(self, denoising_strength: int=10):
        self.denoising_strength = denoising_strength

        with open("config.json") as file:
            self.config = json.load(file)
    
    def __high_contrast(self, img: np.ndarray, threshold: int=220) -> np.ndarray:
        return np.where(img >= threshold, 255, 0)
    
    def __apply_fft(self, image: np.ndarray) -> np.ndarray:
        f = np.fft.fft2(image)
        f_shift = np.fft.fftshift(f)

        # Create a mask to remove high frequency components along the x-axis
        rows, cols = image.shape
        crow, ccol = rows//2 , cols//2
        h = 5
        f_shift[crow-h:crow+h, :] = 0
        f_shift[:, ccol-h:ccol+h] = 0

        # Inverse Fourier Transform
        f_ishift = np.fft.ifftshift(f_shift)
        img_back = np.fft.ifft2(f_ishift)
        img_back = np.abs(img_back)
        new_im = 255 - img_back
        return new_im        
    
    def __apply_denoising(self, img: np.ndarray) -> np.ndarray:
        denoised = cv2.fastNlMeansDenoising(img, None, self.denoising_strength)
        return denoised
    
    def __get_integral(self, x: np.ndarray, y: np.ndarray) -> float:
        return scipy.integrate.simpson(y=y, x=x)

    def __get_integral_percentile(self, x: np.ndarray, y: np.ndarray, percentile: int) -> int:
        full_integral = self.__get_integral(x, y)
        for i in range(len(x) - 1, 0, -1):
            curr_integral = self.__get_integral(x[:i], y[:i])
            if curr_integral / full_integral * 100 < percentile:
                print(x[i-1])
                return x[i - 1]
    
    def __get_contrast_threshold(self, img: np.ndarray, **argv) -> int:
        threshold = [i for i in range(0, 256)]
        sums = []

        for thr in threshold:
            sums.append(self.__high_contrast(img, thr).sum())
        
        return self.__get_integral_percentile(threshold, sums, **argv)
    
    def apply_fft(self, image: np.ndarray) -> np.ndarray:
        return self.__apply_fft(image)
    
    def process_array(self, image: np.ndarray, apply_fft: bool=False) -> np.ndarray:
        if apply_fft:
            image = self.__apply_fft(image)

        return image
    
    def process_image_automatically(self, 
                                    img_path_source: Union[str, np.ndarray], 
                                    apply_fft: bool=False, 
                                    img_path_dest: str=None, **argv) -> np.ndarray:
        
        if isinstance(img_path_source, str):
            img = cv2.imread(img_path_source, 0)
        else:
            img = img_path_source
        
        img = self.process_array(img, apply_fft=apply_fft)
        
        contrast_threshold = self.__get_contrast_threshold(img, **argv)
        processed_img = self.__high_contrast(img, contrast_threshold)

        if img_path_dest:
            cv2.imwrite(img_path_dest, processed_img)
        
        return processed_img, contrast_threshold

    def process_image(self, 
                      img_path_source: Union[str, np.ndarray],
                      apply_fft: bool=False, 
                      contrast_threshold: int=170, 
                      img_path_dest: str=None) -> np.ndarray:
        
        if isinstance(img_path_source, str):
            img = cv2.imread(img_path_source, 0)
        else:
            img = img_path_source

        img = self.process_array(img, apply_fft=apply_fft)
        processed_img = self.__high_contrast(img, contrast_threshold)

        if img_path_dest:
            cv2.imwrite(img_path_dest, processed_img)

        return processed_img
