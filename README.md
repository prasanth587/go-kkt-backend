# mgt

# Build the image and push
```
export DB_CONNECTION_CONN=root:3nRrF3rn5f@tcp(0.0.0.0:3306)/transport_hub
export EMPLOYEE_IMAGE_DIRECTORY=/t_hub_document/employee
export BASE_DIRECTORY=/Users/viadmin/Downloads/mgt
export IMAGE_DIRECTORY=/t_hub_document

Run:
  ./build.sh
  

VERSION=2025.2.19 && docker buildx build -f Dockerfile --platform linux/amd64 -t prabha303/mgthub:$VERSION . && docker push prabha303/mgthub:$VERSION
```
// my name is prahsant 
# IP: 54.172.163.73

# LOGIN API
```
http://<IP>:9005/v1/adminLogin
{
    "message": "User logged in successfully",
    "email_id": "prabhasjraj@gmail.com",
    "first_name": "MGT Transport",
    "login_id": 1,
    "mobile_no": "9876543210",
    "role_id": 1,
    "employee_id": 0,
    "login_type": "Web",
    "last_login": "2025-02-15T16:54:59.413140618Z",
    "organisation": {
        "org_id": 1,
        "name": "MGT Transport",
        "display_name": "MGT Transport",
        "domain_name": "mgt.in",
        "email_id": "prabhasjraj@gmail.com",
        "contact_name": "MGT Admin",
        "contact_no": "9876543210",
        "is_active": 1,
        "logo_path": "",
        "address_line1": "",
        "address_line2": "",
        "city": "Chennai"
    }
}

```

## Create Role:

```
http://<IP>:9005/v1/create/role
{
    "role_code" : "prst_acc",
    "role_name" : "prst accountant",
    "org_id": 1,
    "description1": "prst accountant for finance people"
}
```

# List Role:

 ```
 http://<IP>:9005/v1/1/roles?limit=20&offset=0
[
    {
        "role_id": 3,
        "role_code": "MGT_ACC",
        "role_name": "mgt accountant_5",
        "description": "mgt accountant for finance people",
        "org_id": 1,
        "is_active": 1,
        "version": 9,
        "updated_at": "2025-02-16 12:07:33"
    },
    {
        "role_id": 7,
        "role_code": "PRST_ACC",
        "role_name": "prst accountant",
        "description": "",
        "org_id": 1,
        "is_active": 1,
        "version": 1,
        "updated_at": "2025-02-15 15:01:17"
    },
    {
        "role_id": 5,
        "role_code": "KTT_ACC",
        "role_name": "kkt accountant",
        "description": "kkt accountant for finance people",
        "org_id": 1,
        "is_active": 1,
        "version": 1,
        "updated_at": "2025-02-15 14:47:28"
    },
    ........................................
]


 ```

# Update Role

```
PUT: http://<IP>>:9005/v1/update/3/role
{
    "role_code" : "MGT_ACC",
    "role_name" : "mgt accountant2",
    "org_id": 1,
    "description": "mgt accountant for finance people",
    "is_active": 1
}

```


# create employee and enable login access for the employee

```
POST http://<IP>:9005/v1/create/employee
{
  "first_name": "Alice",
  "last_name": "Johnson",
  "employee_code": "EJ12345",
  "mobile_no": "7829332005",
  "email_id": "alice.johnson@example.com",
  "role_id": 3,
  "dob": "1991-07-06", 
  "gender": "MALE",
  "aadhar_no": "123456789012",
  "access_no": "A1B2C3D4",
  "is_active": 1,
  "joining_date": "2025-03-16", 
  "address_line1": "123 Main St",
  "address_line2": "Apt 4B",
  "city": "Anytown",
  "state": "CA",
  "country": "USA",
  "is_super_admin": 0,
  "is_admin": 0,
  "org_id": 1,
  "login_type": "web",
  "image": "https://example.com/images/alice.jpg"
}

```

# Employee list

```
http://54.91.203.208:9005/v1/create/employee

[
    {
        "emp_id": 2,
        "first_name": "Anderson",
        "last_name": "Johnson",
        "employee_code": "MGT0009",
        "mobile_no": "7829330000",
        "email_id": "anderson.johnson@example.com",
        "role_id": 3,
        "dob": "1991-07-06",
        "gender": "MALE",
        "aadhar_no": "993456789000",
        "access_no": "A1B2C309",
        "is_active": 1,
        "joining_date": "2025-03-16",
        "relieving_date": "",
        "address_line1": "123 Main St",
        "address_line2": "Apt 4B",
        "city": "Anytown",
        "state": "CA",
        "country": "USA",
        "is_super_admin": 0,
        "is_admin": 1,
        "org_id": 1,
        "login_type": "web",
        "image": "https://example.com/images/alice.jpg"
    },
    {
        "emp_id": 1,
        ....
        
    }
]

```


## Employee Active & InActive

```
    PUT http://localhost:9005//v1/employees/14/update/status?isActive=0
    http://localhost:9005//v1/employees/14/update/status?isActive=1
    emp_id = 14
    To Activate = 1
    InActivate = 0
Response:
{
    "name": "Anderson cory",
    "message": "Employee updated successfully",
    "is_active": 0,
    "emp_id": 16
}
```

## Employee Update

 ```
POST http://localhost:9005/v1/employee/16/update

 {
    "first_name": "Anderson cory",
    "last_name": "Johnson cory",
    "employee_code": "MGT001003",
    "mobile_no": "7829330021",
    "email_id": "cryeanderson.johnson@example.com",
    "role_id": 3,
    "dob": "1999-07-06",
    "gender": "MALE",
    "aadhar_no": "993456789000",
    "access_no": "A1B2C3091",
    "is_active": 1,
    "joining_date": "2025-03-16",
    "relieving_date": "2025-03-16",
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "New town",
    "state": "IN",
    "country": "IN",
    "is_super_admin": 1,
    "is_admin": 1,
    "login_type": "web",
    "pin_code": 606804
}   

Response:
{
    "message": "Employee updated successfully: Anderson cory",
    "name": "Anderson cory",
    "emp_id": 16
}

```

# Driver upload API's

```
http://localhost:9005//v1/driver/7/upload/profile?imageFor=profile
http://localhost:9005//v1/driver/7/upload/profile?imageFor=license_front
http://localhost:9005//v1/driver/7/upload/profile?imageFor=license_back
http://localhost:9005//v1/driver/7/upload/profile?imageFor=other

```



# vendor upload 

```

Image for 
 * visiting_card_image
 * pancard_img
 * aadhar_card_img
 * cancelled_check_book_img
 * bank_passbook_img
 * gst_document_img

http://localhost:9005/v1/vendor/8/upload?imageFor=visiting_card_image
http://localhost:9005/v1/vendor/8/upload?imageFor=pancard_img
http://localhost:9005/v1/vendor/8/upload?imageFor=aadhar_card_img
http://localhost:9005/v1/vendor/8/upload?imageFor=cancelled_check_book_img
http://localhost:9005/v1/vendor/8/upload?imageFor=bank_passbook_img
http://localhost:9005/v1/vendor/8/upload?imageFor=gst_document_img

```

# Vehicle upload 
imageFor should be vehicle_image/insurance/registration

```
Image for 
 * vehicle_image
 * fitness_certificate
 * insurance_certificate
 * pollution_certificate
 * national_permits_certificate
 * registration_certificate
 * annual_maintenance_certificate


http://54.227.230.34:9005/v1/vehicle/4/upload?imageFor=vehicle_image
http://54.227.230.34:9005/v1/vehicle/4/upload?imageFor=fitness_certificate
http://54.227.230.34:9005/v1/vehicle/4/upload?imageFor=insurance_certificate

so on

```


#### Customer

# Create customer

```
POST
http://54.175.147.136:9005/v1/create/customer
{
    "customer_name": "Virta Enterprises",
    "customer_code": "C004",
    "credibility_days": 22,
    "company_type": "Public",
    "contact_person": "Michael Davis",
    "mobile_number": "8878889999",
    "alternative_number": "6667778888",
    "email_id": "es.davis@deltaenterprises.com",
    "branch": "Midwest Sol",
    "address": "321 Maple St, Detroit",
    "city": "Detroit",
    "state": "MI",
    "is_active": 1,
    "org_id": 1
}

```

# List Customer

```
GET
http://54.175.147.136:9005/v1/1/customers?limit=20&offset=0

```

# Customer upload image

```
imageFor: agreement_doc_image
POST
http://54.175.147.136:9005/v1/customer/7/upload?imageFor=agreement_doc_image

```

# update customer

```
POST
http://54.175.147.136:9005/v1/customer/5/update

{
    "customer_name": "Delta Enterprises pvt",
    "customer_code": "C0045",
    "credibility_days": 8,
    "company_type": "Public",
    "contact_person": "Michael Davis",
    "mobile_number": "9078889999",
    "alternative_number": "6667778888",
    "email_id": "michael.davis@deltaenterprises.com",
    "branch": "Midwest Branch",
    "address": "321 Maple St, Detroit",
    "city": "Detroit",
    "state": "MI",
    "is_active": 1,
    "org_id": 1
}


```

# Customer active/inactive

```
PUT
http://54.175.147.136:9005/v1/customer/2/update/status?isActive=0

```
