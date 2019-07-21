GC = go

TARGET = learning-bot

all: clean format test $(TARGET)

$(TARGET): main.go
	$(GC) build

format:
	$(GC) fmt ./...

test:
	$(GC) test ./...
	$(GC) vet ./...

clean:
	$(RM) $(TARGET)
