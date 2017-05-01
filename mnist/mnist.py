from tensorflow.examples.tutorials.mnist import input_data
import tensorflow as tf


def max_pooling(inputs):
    return tf.nn.max_pool(
        inputs,
        ksize=[1, 2, 2, 1],
        strides=[1, 2, 2, 1],
        padding='SAME'
    )


def conv2d(inputs, weights):
    return tf.nn.conv2d(
        inputs,
        weights,
        strides=[1, 1, 1, 1],
        padding='SAME'
    )


def create_biases(shape=[]):
    initial = tf.constant(0.1, shape=shape)
    return tf.Variable(initial)


def create_weights(shape=[]):
    initial = tf.truncated_normal(shape, stddev=0.1)
    return tf.Variable(initial)


def create_inputs():
    inputs = tf.placeholder(tf.float32, shape=[None, 784])
    targets = tf.placeholder(tf.float32, shape=[None, 10])
    return inputs, targets


def read_mnist_dataset():
    return input_data.read_data_sets('MNIST_data', one_hot=True)


def main():
    mnist_data = read_mnist_dataset()
    session = tf.InteractiveSession()

    # Create inputs
    inputs, targets = create_inputs()
    # Reshape inputs to a 4d tensor (not sure why)
    x_image = tf.reshape(inputs, [-1, 28, 28, 1])

    # # First convolutional layer
    weights_conv1 = create_weights(shape=[5, 5, 1, 32])
    biases_conv1 = create_biases(shape=[32])

    conv1 = tf.nn.relu(conv2d(x_image, weights_conv1) + biases_conv1)
    conv1 = max_pooling(conv1)

    # # Second convolutional layer
    weights_conv2 = create_weights(shape=[5, 5, 32, 64])
    biases_conv2 = create_biases(shape=[64])

    conv2 = tf.nn.relu(conv2d(conv1, weights_conv2) + biases_conv2)
    conv2 = max_pooling(conv2)

    # #  First fully connected layer
    # idk why 7 * 7
    weigths_fully_con1 = create_weights(shape=[7 * 7 * 64, 1024])
    biases_fully_con1 = create_biases(shape=[1024])
    conv2 = tf.reshape(conv2, [-1, 7 * 7 * 64])
    fully_con1 = tf.nn.relu(tf.matmul(
        conv2, weigths_fully_con1) + biases_fully_con1)

    # # Apply dropout to reduce overfitting
    keep_prob = tf.placeholder(tf.float32)
    dropout = tf.nn.dropout(fully_con1, keep_prob)

    # #  Second fully connected layer
    weigths_fully_con2 = create_weights(shape=[1024, 10])
    biases_fully_con2 = create_biases(shape=[10])
    outputs = tf.matmul(dropout, weigths_fully_con2) + biases_fully_con2

    # # Training
    cross_entropy = tf.nn.softmax_cross_entropy_with_logits(
        labels=targets,
        logits=outputs
    )
    loss = tf.reduce_mean(cross_entropy)

    optimizer = tf.train.AdamOptimizer(1e-4)
    training = optimizer.minimize(loss)

    prediction = tf.equal(tf.argmax(outputs, 1), tf.argmax(targets, 1))
    accuracy = tf.reduce_mean(tf.cast(prediction, tf.float32))

    # # Initialize variables
    session.run(tf.global_variables_initializer())

    for i in range(20000):
        batch = mnist_data.train.next_batch(50)
        _inputs, _targets = batch[0], batch[1]
        if i % 100 == 0:
            train_accuracy = accuracy.eval(
                feed_dict={
                    inputs: _inputs,
                    targets: _targets,
                    keep_prob: 1.0,
                }
            )
            print('Step %d, training accuracy %g' % (i, train_accuracy))
        training.run(feed_dict={
            inputs: _inputs, targets: _targets, keep_prob: 0.5
        })

    print('Test accuracy %g' % accuracy.eval(feed_dict={
        inputs: mnist_data.test.images,
        targets: mnist_data.test.labels,
        keep_prob: 1.0
    }))
if __name__ == '__main__':
    main()
