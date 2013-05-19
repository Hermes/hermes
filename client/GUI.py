from Tkinter import Tk, Text, BOTH, W, N, E, S
from ttk import Frame, Button, Label, Style
from tkFileDialog import *
import sys

class Example(Frame):
  
    def __init__(self, parent):
        Frame.__init__(self, parent)   
         
        self.parent = parent
        
        self.initUI()
        
    def initUI(self):
      
        self.parent.title("Hermes")
        self.style = Style()
        self.style.theme_use("default")
        self.pack(fill=BOTH, expand=1)

        self.columnconfigure(1, weight=1)
        self.columnconfigure(3, pad=7)
        self.rowconfigure(3, weight=1)
        self.rowconfigure(5, pad=7)

        cbtn = Button(self, text="Close", command=self.quit)
        cbtn.grid(row=4, column=3, pady=4)

def main():

    root = Tk()
    root.geometry("350x300+300+300")
    app = Example(root)
    root.mainloop()
    root.destroy()


if __name__ == '__main__':
    main()
