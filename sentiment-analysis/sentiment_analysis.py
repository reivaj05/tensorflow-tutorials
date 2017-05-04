import numpy as np
import tensorflow as tf


def read_file(file_path):
    file = open(file_path, 'r')
    return file.read()


def main():
    reviews = read_file('reviews.txt')
    labels = read_file('labels.txt')
    print labels


if __name__ == '__main__':
    main()
