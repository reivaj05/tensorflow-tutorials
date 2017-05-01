from tensorflow.examples.tutorials.mnist import input_data
import tensorflow as tf


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
    # Create weigths
    weights = tf.Variable(tf.zeros([784, 10]))
    # Create biases
    biases = tf.Variable(tf.zeros([10]))

    # Initializa variables
    session.run(tf.global_variables_initializer())

    # Implement regression model
    outputs = tf.matmul(inputs, weights) + biases

    # Get loss
    cross_entropy = tf.nn.softmax_cross_entropy_with_logits(
        labels=targets, logits=outputs)
    loss = tf.reduce_mean(cross_entropy)

    # Training step
    optimizer = tf.train.GradientDescentOptimizer(0.5)
    training = optimizer.minimize(loss)
    for _ in range(1000):
        batch = mnist_data.train.next_batch(100)
        session.run(training, feed_dict={inputs: batch[0], targets: batch[1]})

    # Check prediction
    prediction = tf.equal(tf.argmax(outputs, 1), tf.argmax(targets, 1))

    # Get accuracy
    accuracy = tf.reduce_mean(tf.cast(prediction, tf.float32))

    print(accuracy.eval(feed_dict={inputs: mnist_data.test.images, targets: mnist_data.test.labels}))


if __name__ == '__main__':
    main()
