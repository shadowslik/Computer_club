package ru.kolobanov.pc.club.exeptions;

public class InvalidHoursException extends RuntimeException {
    public InvalidHoursException(String message) {
        super(message);
    }
}
