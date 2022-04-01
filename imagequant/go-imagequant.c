#include <stdio.h>
#include <stdlib.h>
#include "libimagequant.h"
#include "go-imagequant.h"

imgQuant doQuant(unsigned char* raw_rgba_pixels, unsigned int width, unsigned int height, double gamma) {

    // Use libimagequant to make a palette for the RGBA pixels
    imgQuant workData;
    fprintf(stderr, "Quantization started\n");
    fprintf(stderr, "w: %d, h: %d\n", width, height);
    fprintf(stderr, "dumping head:\n");

    for (unsigned int i = 0; i < 20; i++) {
         fprintf(stderr, "%02X ", raw_rgba_pixels[i]);
    }

    fprintf(stderr, "\n---\n");

    liq_attr *handle = liq_attr_create();
    liq_image *input_image = liq_image_create_rgba(handle, raw_rgba_pixels, (int) width, (int) height, gamma);

    fprintf(stderr, "liq_image_create_rgba() done...\n");

    liq_result *quantization_result;
    if (liq_image_quantize(input_image, handle, &quantization_result) != LIQ_OK) {
        fprintf(stderr, "Quantization failed\n");
        workData.status = EXIT_FAILURE;
        return workData;
    }

    fprintf(stderr, "liq_image_quantize() done...\n");
    // Use libimagequant to make new image pixels from the palette
    // size_t pixels_size = width * height;
    // workData.height = height;
    // workData.width = width;
    workData.size = width * height;
    workData.pixels = malloc(workData.size);

    liq_set_dithering_level(quantization_result, (float) 1.0);

    liq_write_remapped_image(quantization_result, input_image, workData.pixels, workData.size);
    fprintf(stderr, "liq_write_remapped_image() done...\n");

    workData.palette = liq_get_palette(quantization_result);
    fprintf(stderr, "liq_get_palette() done...\n");

    fprintf(stderr, "trying to return workData ...\n");
    return workData;
}

void destroyImgQuant(imgQuant workData) {
    // free(workData.palette);
    free(workData.pixels);
    // free(workData.Palette);
}

