from tkinter import Label, filedialog, Tk, Button

import numpy as np
from PIL import ImageTk, Image
from PIL.Image import Resampling

from face_swaps.face_detect import face_marked, face_marked2
from face_swaps.face_swap import readLandmarkPoints, connectLandmarkPoints, showLandmarkPoints


def open_img():
    x = openfilename()

    img = Image.open(x)

    converted_img = np.array(img)
    showLandmarkPoints(converted_img)
    face_marked_img = face_marked2(converted_img)

    img = Image.fromarray(np.uint8(face_marked_img)).convert('RGB')

    img = img.resize((250, 250), Resampling.LANCZOS)

    img = ImageTk.PhotoImage(img)

    panel = Label(root, image=img)

    panel.image = img
    panel.grid(row=2)


def openfilename():
    filename = filedialog.askopenfilename(title='"pen')
    return filename


root = Tk()

root.title("Image Loader")

root.geometry("550x300+300+150")

root.resizable(width=True, height=True)

Button(root, text='open image', command=open_img).grid(row=1, columnspan=4)
root.mainloop()