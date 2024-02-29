/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("_pb_users_auth_")

  collection.viewRule = "id = @request.auth.id || (id ?~ @collection.blueprints.users.id && @request.auth.id ?~ @collection.blueprints.owner.id )"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("_pb_users_auth_")

  collection.viewRule = "id = @request.auth.id || id ?~ @collection.blueprints.users.id "

  return dao.saveCollection(collection)
})
