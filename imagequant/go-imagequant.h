typedef struct imgQuant {
    int status;
    size_t size;
    const liq_palette *palette;
    unsigned char *pixels;
} imgQuant ;

imgQuant doQuant(unsigned char* raw_rgba_pixels, unsigned int width, unsigned int height, double gamma);
void destroyImgQuant(imgQuant workData);
LIQ_EXPORT void liq_attr_destroy_wrapper(liq_attr *attr) LIQ_NONNULL;
LIQ_EXPORT LIQ_USERESULT liq_image *liq_image_create_rgba_wrapper(const liq_attr *attr, unsigned char *bitmap, int width, int height, double gamma) LIQ_NONNULL;