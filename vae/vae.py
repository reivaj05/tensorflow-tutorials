import functools

from functional import compose, partial

import numpy as np
import tensorflow as tf


class FullyConnected():
    def __init__(self, scope='fully_connected', size=None,
                 dropout=1.0, non_linear=tf.identity):

        assert size, 'Must specify layer size (num nodes)'

        self.scope = scope
        self.size = size
        self.dropout = dropout
        self.non_linear = non_linear

    def __call__(self, _input):
        """apply layer to any input tensor `input`"""
        with tf.name_scope(self.scope):
            while True:
                try:
                    # Reuse weights if already initialized
                    return self.non_linear(
                        tf.matmul(_input, self.weights) + self.bias)
                except(AttributeError):
                    self.initWeightsAndBias(_input)

    def initWeightsAndBias(self, _input):
        self.weights, self.bias = self.wbVars(
            _input.get_shape()[1].value, self.size)

        self.weights = tf.nn.dropout(self.weights, self.dropout)

    @staticmethod
    def wbVars(fan_in, fan_out):
        """Helper to initialize weights and biases, via He's adaptation
        of Xavier init for ReLUs: https://arxiv.org/abs/1502.01852
        """
        # (int, int) -> (tf.Variable, tf.Variable)
        stddev = tf.cast((2 / fan_in)**0.5, tf.float32)

        initial_w = tf.random_normal([fan_in, fan_out], stddev=stddev)
        initial_b = tf.zeros([fan_out])

        return (tf.Variable(initial_w, trainable=True, name="weights"),
                tf.Variable(initial_b, trainable=True, name="biases"))
