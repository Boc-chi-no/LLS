// noinspection JSUnresolvedFunction,JSUnresolvedVariable

db.createUser({
    user: 'shortener',
    pwd: 'VFSNnSFLvfOwFnBh' ,
    roles: [{
        role: 'readWrite',
        db: 'shortener'
    }]
})