from collections import Counter
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
    return reviews.split('\n')


def create_reviews_sequences(reviews, vocab_to_int):
    sequences = []
    for review in reviews:
        sequences.append(
            [vocab_to_int[word] for word in review.split()]
        )
    return sequences


def main():
    # # Read and preprocess data
    reviews = cleanup_reviews(read_file('reviews.txt'))
    labels = read_file('labels.txt')

    text = ' '.join(reviews)
    words = text.split()

    # # Create vocab
    # Count occurrences for each word
    words_counter = Counter(words)
    # Sort the vocab from most common word to least common
    vocab = sorted(words_counter, key=words_counter.get, reverse=True)
    vocab_to_int = {word: index for index, word in enumerate(vocab, 1)}

    # # Create reviews sequences
    reviews_sequences = create_reviews_sequences(reviews, vocab_to_int)


if __name__ == '__main__':
    main()
