from text_filtrator import TextFiltrator

filter = TextFiltrator(contrast_threshold=150, fft_threshold=20)
filter.process_image("test4.jpg", "new_photo.jpg")