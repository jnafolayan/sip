# Sip
Image compression using wavelet transforms.

**Wavelet families implemented**
- Haar
- Cohen-Daubechies-Feauveau

## Installation
```sh
git clone https://github.com/jnafolayan/sip
cd sip

# Compress an image
# wavelet types: haar, cdf97
# level: 1, 2, 3, ...
# threshold: 0, 1, 2, ...
go run main.go compress PATH_TO_FILE -wavelet WAVELET_TYPE -level LEVEL -output OUTPUT_FILE -threshold THRESHOLD

```