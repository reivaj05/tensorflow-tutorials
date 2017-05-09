from collections import Counter
import numpy as np
from string import punctuation
import tensorflow as tf


def get_batches(x, y, batch_size=100):
    num_batches = len(x) // batch_size
    x, y = x[:num_batches * batch_size], y[:num_batches * batch_size]
    for i in range(0, len(x), batch_size):
        yield x[i: i + batch_size], y[i: i + batch_size]


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

    # # Embedding
    embedding_size = 300
    with graph.as_default():
        embedding = tf.Variable(
            tf.random_uniform((num_words, embedding_size), -1, 1))
        embed = tf.nn.embedding_lookup(embedding, _inputs)

    # # Create lstm cell

    with graph.as_default():
        lstm = tf.contrib.rnn.BasicLSTMCell(num_units)
        drop = tf.contrib.rnn.DropoutWrapper(lstm, output_keep_prob=keep_prob)
        cell = tf.contrib.rnn.MultiRNNCell([drop] * lstm_layers)

        initial_state = cell.zero_state(batch_size, tf.float32)

    # # Forwarding

    with graph.as_default():
        outputs, final_state = tf.nn.dynamic_rnn(
            cell, embed, initial_state=initial_state)

    # # Output. Fully connected layer

    with graph.as_default():
        predictions = tf.contrib.layers.fully_connected(
            outputs[:, -1],
            1,
            activation_fn=tf.sigmoid
        )
        error = tf.losses.mean_squared_error(_targets, predictions)
        optimizer = tf.train.AdamOptimizer(learning_rate).minimize(error)

    # # Validation accuracy

    with grap.as_default():
        expected_prediction = tf.equal(
            tf.cast(tf.round(predictions), tf.int32), _targets)
        accuracy = tf.reduce_mean(tf.cast(expected_prediction, tf.float32))


if __name__ == '__main__':
    main()
