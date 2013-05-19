from Tkinter import Tk, Text, BOTH, W, N, E, S
from ttk import Frame, Button, Label, Style
from TkFileDialog import *
import subprocess


class Example(Frame):
  
    def __init__(self, parent):
        Frame.__init__(self, parent)   
         
        self.parent = parent
        
        self.initUI()
        
    def initUI(self):
      
        self.parent.title("Windows")
        self.style = Style()
        self.style.theme_use("default")
        self.pack(fill=BOTH, expand=1)

        self.columnconfigure(1, weight=1)
        self.columnconfigure(3, pad=7)
        self.rowconfigure(3, weight=1)
        self.rowconfigure(5, pad=7)
        
        lbl = Label(self, text="Windows")
        lbl.grid(sticky=W, pady=4, padx=5)
        
        area = Text(self)
        area.grid(row=1, column=0, columnspan=2, rowspan=4, 
            padx=5, sticky=E+W+S+N)

        p1btn = Button(self, text="Push")
        p1btn.grid(row=1, column=3)

        p2btn = Button(self, text="Pull")
        p2btn.grid(row=2, column=3, pady=4)

        ubtn = Button(self, text="Upload a folder", command=TkFileDialog.askdirectory())
        ubtn.grid(row=2, column=3, pady=4)

        cbtn = Button(self, text="Close", command=self.quit)
        cbtn.grid(row=3, column=4, pady=4)


def main():
  
    root = Tk()
    root.geometry("350x300+300+300")
    app = Example(root)
    root.mainloop()  


if __name__ == '__main__':
    main() 
