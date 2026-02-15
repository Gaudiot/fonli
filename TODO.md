# TODO - Italian Learning Platform

This document serves as a TDD (Test-Driven Development) guide for the Italian learning platform development.

## v0.1.0: Environment Configuration

### Tasks
- [X] Set up `.env` file structure
- [X] Configure environment variables for API keys
- [X] Add `.env.example` file with template
- [X] Implement environment variable validation

## v0.2.0: Portuguese to Italian Translation Exercise

### Tasks
- [X] Create exercise endpoint/service for Portuguese to Italian translation
- [X] Implement AI integration for generating 10 simple Portuguese words
- [X] Implement AI integration for generating correct Italian translations
- [X] Create API endpoint that returns words with their correct translations

## v0.3.0: Italian to Portuguese Translation Exercise

### Tasks
- [X] Create exercise endpoint/service for Italian to Portuguese translation
- [X] Implement AI integration for generating 10 simple Italian words
- [X] Implement AI integration for generating correct Portuguese translations
- [X] Create API endpoint that returns words with their correct translations

## v0.4.0: Italian Verb Conjugation Exercise

### Tasks
- [X] Create verb conjugation endpoint/service
- [X] Implement verb selection mechanism
- [X] Add tense selection feature
- [X] Implement AI integration for generating verb conjugations in all persons (1st, 2nd, 3rd - singular and plural)
- [X] Create API endpoint that returns verb with all conjugations for selected tense

## v0.5.0: Short Story Translation Exercise

### Tasks
- [ ] Create story translation endpoint/service
- [ ] Implement AI integration for generating short Portuguese stories
- [ ] Create API endpoint to receive user's translation
- [ ] Implement AI integration for correcting and providing feedback on user's translation
- [ ] Create API endpoint that returns corrected translation with feedback

## General TDD Guidelines

- Write tests before implementing features
- Ensure all tests pass before moving to the next task
- Refactor code while keeping tests green
- Maintain test coverage above 80%
- Use descriptive test names that explain the behavior being tested
