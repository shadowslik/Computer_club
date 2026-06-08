package ru.kolobanov.pc.club.exeptions;

public class ComputerTypeNotFoundException extends RuntimeException {
    public ComputerTypeNotFoundException(String message) {
        super(message);
    }
}
