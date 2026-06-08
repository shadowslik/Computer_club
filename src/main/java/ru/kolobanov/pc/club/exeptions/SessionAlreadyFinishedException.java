package ru.kolobanov.pc.club.exeptions;

public class SessionAlreadyFinishedException extends RuntimeException {
    public SessionAlreadyFinishedException(String message) {
        super(message);
    }
}
