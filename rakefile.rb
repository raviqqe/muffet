task :deps do
  sh 'go get -u github.com/alecthomas/gometalinter'
  sh 'gometalinter --install'
  sh 'go get -d -t ./...'
end

task :lint do
  sh 'gometalinter --disable vetshadow ./...'
end

task :build do
  sh 'go build'
end

task :unit_test do
  sh 'go test -covermode atomic -coverprofile coverage.txt'
end

task integration_test: :build do
  sh './muffet http://localhost:8080'
  sh '! ./muffet http://localhost:8888'

  sh './muffet -c 1 http://localhost:8080'
  sh './muffet --concurrency 1 http://localhost:8080'

  sh './muffet -n 1 http://localhost:8080'
  sh './muffet --connections-per-host 1 http://localhost:8080'

  sh './muffet --help'

  sh './muffet -v http://localhost:8080 2>&1 | grep OK'
  sh './muffet --verbose http://localhost:8080 2>&1 | grep OK'
  sh '! ./muffet http://localhost:8080 2>&1 | grep OK'

  sh './muffet -v http://localhost:8080 2>&1 | sort > /tmp/muffet_1.txt'
  sh './muffet -v http://localhost:8080 2>&1 | sort > /tmp/muffet_2.txt'
  sh 'diff /tmp/muffet_1.txt /tmp/muffet_2.txt'

  sh '! ./muffet http://localhost:8080 2>&1 | grep .'
end

task :serve do
  cd 'test/valid' do
    sh 'python3 -m http.server 8080 &'
  end

  cd 'test/dead_link' do
    sh 'python3 -m http.server 8888 &'
  end
end
