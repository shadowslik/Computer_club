package ru.kolobanov.pc.club.exeptions;

public class PhoneAlreadyExistException extends RuntimeException {
    public PhoneAlreadyExistException(String message) {
        super(message);
    }
}
