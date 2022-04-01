#pragma pack(0)
typedef struct imgQuant {
    int status;
    unsigned char *pixels;
    size_t size;
    const liq_palette *palette;
} imgQuant ;

imgQuant doQuant(unsigned char* raw_rgba_pixels, unsigned int width, unsigned int height, double gamma);
void destroyImgQuant(imgQuant workData);