import numpy as np
from string import punctuation
import tensorflow as tf


def read_file(file_path):
    file = open(file_path, 'r')
    return file.read()


def cleanup_reviews(reviews):
    reviews = ''.join(
        [char for char in reviews if char not in punctuation]
    )
    reviews = reviews.split('\n')
    return ' '.join(reviews)


def main():
    reviews = cleanup_reviews(read_file('reviews.txt'))
    labels = read_file('labels.txt')
    print reviews[:2000]



if __name__ == '__main__':
    main()
