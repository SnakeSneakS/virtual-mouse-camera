# VIRTUAL-MOUSE-CAMERA
- virtual mouse controlled by hand gesture 


# ARCHITECTURE
camera -> hand-tracker -> UDP SOCKET (hand landmarks) -> mouse-controller -> manipulate mouse

# setup & run
- see [hand-tracker](./hand-tracker/README.md) and [mouse-controller](./mouse-controller/README.md)

# run using makefile
- `make run-hand-tracker ADDR=localhost PORT=8080 CAMERA_INDEX=1`
- `make run-mouse-controller ADDR=localhost PORT=8080`

# usage
- run
- move yout hand in front of camera
- touch your index finger by middle finger to hold mouse's left. If you want to click, touch and release soon. 
