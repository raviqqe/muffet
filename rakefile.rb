task :deps do
  sh 'go get -u github.com/alecthomas/gometalinter'
  sh 'gometalinter --install'
  sh 'go get -d -t ./...'
end

task :lint do
  options = %w[gas gocyclo gosec maligned vetshadow].map do |l|
    "--disable #{l}"
  end.join ' '

  sh "gometalinter #{options} ./..."
end

task :build do
  sh 'go build -ldflags "-X main.version=$(git describe --abbrev=0 --tags)"'
end

task :unit_test do
  sh 'go test -covermode atomic -coverprofile coverage.txt'
end

task integration_test: :build do
  sh './muffet http://localhost:1111'
  sh '! ./muffet http://localhost:2222'

  sh './muffet -c 1 http://localhost:1111'
  sh './muffet --concurrency 1 http://localhost:1111'

  sh './muffet --help'

  sh './muffet -v http://localhost:1111 | grep 200'
  sh '[ $(./muffet -v http://localhost:1111 | wc -l) -eq 8 ]'
  sh './muffet --verbose http://localhost:1111 | grep 200'
  sh '! ./muffet http://localhost:1111 | grep 200'

  sh './muffet -v http://localhost:1111 | sort > /tmp/muffet_1.txt'
  sh './muffet -v http://localhost:1111 | sort > /tmp/muffet_2.txt'
  sh 'diff /tmp/muffet_1.txt /tmp/muffet_2.txt'

  sh '[ $(./muffet -rv http://localhost:1111 | wc -l) -eq 6 ]'
  sh '[ $(./muffet -sv http://localhost:1111 | wc -l) -eq 6 ]'

  sh '! ./muffet http://localhost:1111 | grep .'
end

task test: %w[unit_test integration_test]

task :serve do
  [['test/valid', 1111], ['test/dead_link', 2222]].each do |args|
    sh "ruby -run -e httpd #{args[0]} -p #{args[1]} &"
  end
end
