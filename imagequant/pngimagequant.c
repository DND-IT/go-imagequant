#include <stdio.h>
#include <stdlib.h>
#include "lodepng.h"
#include "libimagequant.h"
#include "pngimagequant.h"

pngQuant doQuant(unsigned char* raw_rgba_pixels, unsigned int width, unsigned int height, double gamma) {

    // Use libimagequant to make a palette for the RGBA pixels
    pngQuant workData;
    workData.Status = 0;

    fprintf(stderr, "Quantization started\n");
    fprintf(stderr, "w: %d, h: %d\n", width, height);
    // fprintf(stderr, "sizeof raw_rgba_pixel: %d\n",sizeof(raw_rgba_pixels));
    // fprintf(stderr, "ptr raw_rgba_pixel: %p\n",raw_rgba_pixels);
    fprintf(stderr, "dumping head:\n");
    for (unsigned int i = 0; i < 20; i++) {
        fprintf(stderr, "%02X ", raw_rgba_pixels[i]);
    }
    fprintf(stderr, "\n---\n");

    liq_attr *handle = liq_attr_create();
    liq_image *input_image = liq_image_create_rgba(handle, raw_rgba_pixels, width, height, gamma);

    fprintf(stderr, "liq_image_create_rgba() done...\n");

    liq_result *quantization_result;
    if (liq_image_quantize(input_image, handle, &quantization_result) != LIQ_OK) {
        fprintf(stderr, "Quantization failed\n");
        workData.Status = EXIT_FAILURE;
        return workData;
    }

    fprintf(stderr, "liq_image_quantize() done...\n");
    // Use libimagequant to make new image pixels from the palette
    size_t pixels_size = width * height;
    unsigned char *raw_8bit_pixels = malloc(pixels_size);
    liq_set_dithering_level(quantization_result, 1.0);

    liq_write_remapped_image(quantization_result, input_image, raw_8bit_pixels, pixels_size);
    fprintf(stderr, "liq_write_remapped_image() done...\n");

    const liq_palette *palette = liq_get_palette(quantization_result);
    fprintf(stderr, "liq_get_palette() done...\n");


    LodePNGState state;
    lodepng_state_init(&state);
    state.info_raw.colortype = LCT_PALETTE;
    state.info_raw.bitdepth = 8;
    state.info_png.color.colortype = LCT_PALETTE;
    state.info_png.color.bitdepth = 8;

    for(int i=0; i < palette->count; i++) {
        lodepng_palette_add(&state.info_png.color, palette->entries[i].r, palette->entries[i].g, palette->entries[i].b, palette->entries[i].a);
        lodepng_palette_add(&state.info_raw, palette->entries[i].r, palette->entries[i].g, palette->entries[i].b, palette->entries[i].a);
    }


    fprintf(stderr, "trying to call lodepng_encode() ...\n");
    unsigned int out_status = lodepng_encode(&workData.Png, &workData.Size, raw_8bit_pixels, width, height, &state);
    if (out_status) {
        workData.Status = EXIT_FAILURE;
        return workData;
    }

    fprintf(stderr, "lodepng_encode() done ...\n");
    fprintf(stderr, "workData.Size: %u\n", workData.Size);

    liq_result_destroy(quantization_result); // Must be freed only after you're done using the palette
    liq_image_destroy(input_image);
    liq_attr_destroy(handle);

    free(raw_8bit_pixels);
    lodepng_state_cleanup(&state);
    fprintf(stderr, "trying to return workData ...\n");

    return workData;

}