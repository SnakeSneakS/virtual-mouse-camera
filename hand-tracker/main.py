# %%
import argparse
from turtle import position
import cv2
import mediapipe as mp
from socket import socket, AF_INET, SOCK_DGRAM
import pickle
import json

# %% argument parser 
parser = argparse.ArgumentParser(description='add addr:port')
parser.add_argument('addr', type=str, help='address to send hand info')
parser.add_argument('port', type=int, help='port to send hand info')
parser.add_argument('camera_index', type=int, help='which camera to use')
args = parser.parse_args()
print(args)

# %% mediapipe setup
mp_drawing = mp.solutions.drawing_utils
mp_drawing_styles = mp.solutions.drawing_styles
mp_hands = mp.solutions.hands


# %% camera speed setup
FPS = 10
waitKeyTime = round(1000 / FPS)

# %% capture camera and send hand data to addr:port
def start(addr, port, cameraIndex):
    # socket setup
    s = socket(AF_INET, SOCK_DGRAM)

    # For webcam input:
    cap = cv2.VideoCapture(cameraIndex)

    # print info
    print("addr: ", addr)
    print("port: ", port)
    print("cameraIndex: ", cameraIndex)

    dest = (addr,port)

    # handle camera image
    with mp_hands.Hands(
        model_complexity=0,
        min_detection_confidence=0.5,
        min_tracking_confidence=0.5) as hands:

        while cap.isOpened():
            # camera image 
            success, image = cap.read()
            if not success:
                print("ingoring empty camera frame.")
                continue

            # hands
            results = hands.process(image)
            image.flags.writeable = True
            if results.multi_hand_landmarks:
                for hand_landmarks in results.multi_hand_landmarks:
                    mp_drawing.draw_landmarks(
                        image,
                        hand_landmarks,
                        mp_hands.HAND_CONNECTIONS,
                        mp_drawing_styles.get_default_hand_landmarks_style(),
                        mp_drawing_styles.get_default_hand_connections_style()
                    )

            #show image
            cv2.imshow('MediaPipe Hands', cv2.flip(image, 1))

            # send UDP socket to addr:port
            if results.multi_hand_landmarks:
                hand_landmarks = results.multi_hand_landmarks[0] # 片方の手だけ
                
                # landmarkをlistとして組み立てる
                landmarks = []
                for i in range(21):
                    landmarks.append(
                        {
                            "x": hand_landmarks.landmark[i].x,
                            "y": hand_landmarks.landmark[i].y, 
                            "z": hand_landmarks.landmark[i].z
                        }
                        #hand_landmarks.landmark[i]
                    )
                #print( "landmarks: ", landmarks )
                landmarks_json = json.dumps( landmarks )
                landmarks_byte = landmarks_json.encode("ascii")
                #print(landmarks_byte.decode("utf-8"))
                s.sendto(landmarks_byte, dest) 

            if cv2.waitKey(waitKeyTime) & 0xFF == 27:
                break

    cap.release()

# %% run
start(args.addr, args.port, args.camera_index)