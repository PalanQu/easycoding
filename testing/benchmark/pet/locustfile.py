from locust import HttpUser, task

class PetService(HttpUser):
    @task
    def put_pet(self):
        self.client.get('')

    @task
    def get_pet(self):
        self.client.get('')

    @task
    def delete_pet(self):
        self.client.get('')
