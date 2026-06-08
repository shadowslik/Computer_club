package ru.kolobanov.pc.club.exeptions;

public class BankUnavailableException extends RuntimeException {
    public BankUnavailableException(String message) {
        super(message);
    }
}
