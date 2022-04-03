#include <stdio.h>
#include <stdlib.h>
#include "libimagequant.h"
#include "go-imagequant.h"


void liq_attr_destroy_wrapper(liq_attr *attr) {
    // fprintf(stderr, "calling liq_attr_destroy with ptr to liq_attr: %p\n", attr);
    liq_attr_destroy(attr);
}

// liq_image_create_rgba_wrapper is a wrapper for liq_image_create_rgba.
LIQ_USERESULT liq_image *liq_image_create_rgba_wrapper(const liq_attr *attr, unsigned char *raw_rgba_pixels, int width, int height, double gamma) {
 return liq_image_create_rgba(attr, raw_rgba_pixels, width, height, gamma);
}



