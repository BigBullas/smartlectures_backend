import numpy as np
import cv2
import json


class TextFiltrator():
    def __init__(self, contrast_threshold: int=220, fft_threshold: int=20, 
                 contrash_threshold_for_fft: int=130, denoising_strength: int=10):
        self.contrast_threshold = contrast_threshold
        self.fft_threshold = fft_threshold
        self.contrast_threshold_for_fft = contrash_threshold_for_fft
        self.denoising_strength = denoising_strength

        with open("config.json") as file:
            self.config = json.load(file)
    
    def __rgb2gray(self, rgb: np.ndarray) -> np.ndarray:
        return np.dot(rgb, [0.2989, 0.5870, 0.1140])
    
    def __high_contrast(self, img: np.ndarray, threshold: int=220) -> np.ndarray:
        return np.where(img >= threshold, 255, 0)
    
    def __apply_fft(self, img: np.ndarray) -> np.ndarray:
        img_fft = np.fft.fft2(img)
        mags = np.abs(np.fft.fftshift(img_fft))
        angles = np.angle(np.fft.fftshift(img_fft))
        visual = np.log(mags)
        visual2 = (visual - visual.min()) / (visual.max() - visual.min())*255

        visual2_contrast = self.__high_contrast(visual2, self.contrast_threshold_for_fft)

        mask_y = (np.mean(visual2_contrast, -1) > self.fft_threshold)
        mask_x = (np.mean(visual2_contrast, 0) > self.fft_threshold)
        visual[mask_y] = np.mean(visual)
        visual = visual.T
        visual[mask_x] = np.mean(visual)
        visual = visual.T

        newmagsshift = np.exp(visual)
        newffts = newmagsshift * np.exp(1j*angles)
        newfft = np.fft.ifftshift(newffts)
        imrev = np.fft.ifft2(newfft)
        newim2 = 255 - np.abs(imrev).astype(np.uint8)

        return newim2
    
    def __apply_denoising(self, img: np.ndarray) -> np.ndarray:
        denoised = cv2.fastNlMeansDenoising(img, None, self.denoising_strength)
        return denoised
    
    def process_array(self, img: np.ndarray) -> np.ndarray:
        if self.config["rgb2gray"]:
            img = self.__rgb2gray(img)
        if self.config["high_contrast"]:
            img = self.__high_contrast(img, self.contrast_threshold)
        if self.config["apply_fft"]:
            img = self.__apply_fft(img)
        if self.config["apply_denoising"]:
            img = self.__apply_denoising(img)
            img = self.__high_contrast(img, 240)

        return img

    def process_image(self, img_path_source: str, img_path_dest: str=None) -> None:
        img = cv2.imread(img_path_source)[:,:,:3]
        img = self.process_array(img)
        cv2.imwrite(img_path_dest, img)
