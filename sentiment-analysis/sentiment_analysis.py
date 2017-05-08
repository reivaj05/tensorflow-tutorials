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

    # # Encoding labels positive/negative (Cast them to 0 and 1)

    labels = labels.split('\n')
    labels = np.array([1 if label == 'positive' else 0 for label in labels])

    # # Remove zero length reviews/labels
    non_zero_idx = [index for index, review in enumerate(reviews_sequences)if len(review) != 0]
    reviews_sequences = [reviews_sequences[index] for index in non_zero_idx]
    labels = np.array([labels[index] for index in non_zero_idx])

    # # Create sequences with only first 200 words per review
    sequnce_size = 200
    sequences = np.zeros((len(reviews_sequences), sequnce_size), dtype=int)
    for i, review in enumerate(reviews_sequences):
        sequences[i, -len(review):] = np.array(review)[:sequnce_size]

    # # Split, validation and test
    split_fraction = 0.8
    split_index = int(len(sequences) * split_fraction)

    training_input = sequences[:split_index]
    validation_input = sequences[split_index:]

    training_target = labels[:split_index]
    validation_target = labels[split_index:]

    testing_index = int(len(validation_input) * 0.5)
    testing_input = validation_input[testing_index:]
    validation_input = validation_input[:testing_index]
    testing_target = validation_target[testing_index:]
    validation_target = validation_target[:testing_index]

    # # Build graph
    # Hyperparameters
    num_units = 256
    lstm_layers = 1
    batch_size = 500
    learning_rate = 0.001

    num_words = len(vocab_to_int)
    graph = tf.Graph()
    with graph.as_default():
        _inputs = tf.placeholder(tf.int32, [None, None], name='inputs')
        _targets = tf.placeholder(tf.int32, [None, None], name='targets')
        keep_prob = tf.placeholder(tf.float32, name='keep_prob')

    

if __name__ == '__main__':
    main()
