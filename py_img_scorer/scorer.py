import tkinter as tk
from tkinter import ttk
from tkinter import filedialog
from PIL import Image, ImageTk
import os, sys
from os import path


####################### ImageGenerator #############################
class ImageGenerator:

    def __init__(self, workingDir, targetImgHeight=1000):
        self.workingDir = workingDir
        self.index = 0
        self.imgNameList = list(self.img_gen(workingDir))
        self.currentImg = None
        self.nextImg = None
        self.prevImg = None
        self.targetImgHeight = targetImgHeight

        if len(self.imgNameList) > 0:
            self.init_img()

    def init_img(self):
        self.currentImg = self.img_tk(self.index)
        self.nextImg = self.img_tk(self.next_index(self.index))
        self.prevImg = self.img_tk(self.prev_index(self.index))

    @staticmethod
    def img_gen(wrorkingDir):
        for f in os.listdir(wrorkingDir):
            if f.endswith('.jpg') or f.endswith('.png') or f.endswith(
                    '.jpeg') or f.endswith('.JPG') or f.endswith(
                        '.PNG') or f.endswith('.JPEG'):
                yield f

    @staticmethod
    def resize_img(img, width=-1, height=-1, inter=Image.BICUBIC):
        if width < 0 and height < 0:
            return img
        elif width < 0:
            img_size = img.size
            resize_ratio = height / img_size[1]
            new_size = (round(img_size[0] * resize_ratio),
                        round(img_size[1] * resize_ratio))
            return img.resize(new_size, inter)
        elif height < 0:
            img_size = img.size
            resize_ratio = width / img_size[0]
            new_size = (round(img_size[0] * resize_ratio),
                        round(img_size[1] * resize_ratio))
            return img.resize(new_size, inter)

    def img_tk(self, index):
        img = Image.open(path.join(self.workingDir, self.imgNameList[index]))
        img = self.resize_img(img, height=self.targetImgHeight)
        img = ImageTk.PhotoImage(img)
        return img

    def set_working_dir(self, workingDir):
        self.workingDir = workingDir
        self.index = 0
        self.imgNameList = list(self.img_gen(workingDir))
        self.init_img()

    def next_index(self, index):
        return (index + 1) % len(self.imgNameList)

    def prev_index(self, index):
        return (index - 1) % len(self.imgNameList)

    def next(self):
        self.index = self.next_index(self.index)
        self.prevImg = self.currentImg
        self.currentImg = self.nextImg
        self.nextImg = self.img_tk(self.next_index(self.index))
        return self.currentImg

    def prev(self):
        self.index = self.prev_index(self.index)
        self.nextImg = self.currentImg
        self.currentImg = self.prevImg
        self.prevImg = self.img_tk(self.prev_index(self.index))
        return self.currentImg

    def get_img(self):
        return self.currentImg


############################## APP ###################################
class App(tk.Tk):

    def __init__(self):
        super().__init__()
        self.title("Image Scorer")

        self.workingDir = sys.argv[1] if (
            len(sys.argv) > 1 and os.path.exists(sys.argv[1])) else '.'

        self.imgGenerator = ImageGenerator(self.workingDir)
        self.currImg = self.imgGenerator.get_img()

        self.mainFrame = self.main_frame(self)
        self.mainMenu = self.main_menu(self)
        self.imageFrame = self.img_frame(self.mainFrame)
        self.imgLabel = self.img_label(self.imageFrame)
        self.nextBtn = self.next_btn(self.imageFrame, self.next_action)
        self.prevBtn = self.back_btn(self.imageFrame, self.prev_action)
        self.prevBtn.configure(command=self.prev_action)

        self.update_img()

        self.bind('<Right>', lambda e: self.next_action())
        self.bind('<Left>', lambda e: self.prev_action())
        self.bind('<Home>', lambda e: print('home'))
        self.bind('<Escape>', lambda e: self.quit())

    @staticmethod
    def main_frame(parent):
        mainFrame = ttk.Frame(parent)
        mainFrame.pack(fill=tk.BOTH, expand=True)
        return mainFrame

    @staticmethod
    def img_frame(parent):
        imgFrame = ttk.Frame(parent)
        imgFrame.grid(row=0, column=0, columnspan=3, sticky=tk.NSEW)
        return imgFrame

    @staticmethod
    def img_label(parent):
        imgLabel = ttk.Label(parent)
        imgLabel.grid(row=0, column=1, sticky=tk.NSEW)
        return imgLabel

    @staticmethod
    def next_btn(parent, command):
        nextBtn = ttk.Button(parent, text='>>', command=command)
        nextBtn.grid(row=0, column=2, sticky=tk.NSEW)
        return nextBtn

    @staticmethod
    def back_btn(parent, command):
        backBtn = ttk.Button(parent, text='<<', command=command)
        backBtn.grid(row=0, column=0, sticky=tk.NSEW)
        return backBtn

    @staticmethod
    def main_menu(root):
        mainMenu = tk.Menu(root)

        fileMenu = tk.Menu(mainMenu, tearoff=0)
        fileMenu.add_command(label='Open', command=root.file_open_action)
        fileMenu.add_command(label='Save', command=lambda: print('save'))
        fileMenu.add_separator()
        fileMenu.add_command(label='Exit', command=lambda: root.quit())
        mainMenu.add_cascade(label='File', menu=fileMenu)

        root.config(menu=mainMenu)
        return mainMenu

    def update_img(self):
        self.currImg = self.imgGenerator.get_img()
        if self.currImg:
            self.imgLabel.config(image=self.currImg)
            self.imgLabel.image = self.currImg

    def file_open_action(self):
        self.workingDir = filedialog.askdirectory(initialdir='.',
                                                  mustexist=True)
        self.imgGenerator.set_working_dir(self.workingDir)
        self.currImg = self.imgGenerator.get_img()
        self.update_img()

    def next_action(self):
        self.currImg = self.imgGenerator.next()
        self.update_img()

    def prev_action(self):
        self.currImg = self.imgGenerator.prev()
        self.update_img()


app = App()
app.mainloop()
