package ru.kolobanov.pc.club.exeptions;

public class ComputerTypeAlreadyExistsException extends RuntimeException {
    public ComputerTypeAlreadyExistsException(String message) {
        super(message);
    }
}
