Small Golang project I made to simply pixelate images.

To run, simply take your image path and:
```
go run main.go <image-before>.jpg <image-to-create>.jpg {pixelate count} {number of colors}
```

For the example below, the first image is the original:

![lion](https://github.com/user-attachments/assets/0b1af293-154c-480a-ab98-54d092534e03)

This next one is using the command with `10` as the `pixelate count` and `2` as the `number of colors`:
```
go run main.go lion.jpg lion-pixelate.jpg 10 2
```

![lion-pixelate](https://github.com/user-attachments/assets/1343d697-92a0-4eda-960b-09bbce185ac4)

Then this one has `8` as the `number of colors`:
```
go run main.go lion.jpg lion-pixelate.jpg 10 8
```

![lion-pixelate](https://github.com/user-attachments/assets/98a5762b-399e-468b-9895-19297af5f4a0)

Finally, this has `50` as the `pixleate counte`:
```
go run main.go lion.jpg lion-pixelate.jpg 50 8
```

![lion-pixelate](https://github.com/user-attachments/assets/9c7cfcc7-ae2e-4132-a508-1461155a9bf3)
