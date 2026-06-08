package ru.kolobanov.pc.club.exeptions;

public class ComputerNotFoundException extends RuntimeException {
    public ComputerNotFoundException(String message) {
        super(message);
    }
}
