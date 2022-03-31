#pragma pack(0)
typedef struct {
    int Status;
    unsigned char *Png;
    size_t Size;
} pngQuant ;

pngQuant doQuant(unsigned char* raw_rgba_pixels, unsigned int width, unsigned int height, double gamma);
