echo "Installing Python for boostrapping"
sudo apt install python3

echo "Installing poetry"
curl -sSL https://install.python-poetry.org | python3 -
poetry install

echo "Installing Node"
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

echo "Installing packages"
sudo apt install postgresql built-essential

echo "Installing bazel"
sudo npm install -g @bazel/bazelisk
