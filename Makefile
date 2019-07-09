GC = go

TARGET = learning-bot

all: $(TARGET)

$(TARGET): main.go
	$(GC) build

clean:
	$(RM) $(TARGET)
