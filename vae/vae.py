import numpy as np
import tensorflow as tf
from tensorflow.examples.tutorials.mnist import input_data

import matplotlib.pyplot as plt

np.random.seed(0)
tf.set_random_seed(0)

mnist = input_data.read_data_sets('MNIST_data', one_hot=True)
num_samples = mnist.train.num_examples


def xavier_init_weights(fan_in, fan_out, constant=1):
    """ Xavier initialization of network weights"""
    # https://stackoverflow.com/questions/33640581/how-to-do-xavier-initialization-on-tensorflow
    a = np.sqrt(6.0 / (fan_in + fan_out))
    low, high = -constant * a, constant * a
    return tf.random_uniform((fan_in, fan_out),
                             minval=low, maxval=high,
                             dtype=tf.float32)



# TODO