package ru.kolobanov.pc.club.exeptions;

public class ComputerBusyException extends RuntimeException {
    public ComputerBusyException(String message) {
        super(message);
    }
}
